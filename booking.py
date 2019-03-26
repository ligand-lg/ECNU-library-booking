import requests
import time
import json
import logging
import os

logging.basicConfig(level=logging.INFO,
                    format='%(asctime)-s %(levelname)-s: %(message)s')
# 最大提前天数：2天
MAX_DELAY_DAY = 2


def get_config(delayDay=2):
    """ 获取配置文件,从同目录下的conf.json文件加载 """
    file_path = os.path.join(os.path.dirname(__file__), 'conf.json')
    with open(file_path, encoding="utf-8") as cf:
        conf = json.load(cf)
        conf['delayDay'] = delayDay
        return conf


def check_config(config):
    sid = config.get('sid', None)
    if not (isinstance(sid, str) and len(sid) > 1):
        logging.error('sid: 学号填写有误')
        return False
    password = config.get('password', None)
    if not (isinstance(password, str) and len(password) > 1):
        logging.error('password: 公共数据库密码有误')
        return False
    roomNo = config.get('roomNo', None)
    if not (isinstance(roomNo, str) and len(roomNo) == 4):
        logging.error('roomNo: 房间号有误')
        return False
    startTime = config.get('startTime', None)
    if not (isinstance(startTime, str) and len(startTime) == 5):
        logging.error('startTime: 开始时间有误')
        return False
    endTime = config.get('endTime', None)
    if not (isinstance(endTime, str) and len(endTime) == 5):
        logging.error('endTime: 结束事件有误')
        return False
    delayDay = config.get('delayDay', None)
    if not (isinstance(delayDay, int) and delayDay <= MAX_DELAY_DAY):
        logging.error('delayDay: 预订日期有误，不能大于'+MAX_DELAY_DAY)
        return False
    return True


def get_rooms():
    """ 可以选择的房间. 从当前目录下的 rooms.json 文件中读取 """
    rooms = []
    zhongbei_file_name = 'zhongbei_rooms.json'
    minghang_file_name = 'minghang_rooms.json'
    dir_path = os.path.dirname(__file__)
    with open(os.path.join(dir_path, zhongbei_file_name), encoding="utf-8") as f:
        rooms.extend(json.load(f)['data'])
    with open(os.path.join(dir_path, minghang_file_name), encoding="utf-8") as f:
        rooms.extend(json.load(f)['data'])
    return rooms


class Booking(object):
    def __init__(self, config):
        global MAX_DELAY_DAY
        self.config = config
        # 学校预订系统的相关信息，可能会更新
        self.host = '202.120.82.2'
        self.port = 8081
        self.login_url = '/ClientWeb/pro/ajax/login.aspx'
        self.booking_url = '/ClientWeb/pro/ajax/reserve.aspx'
        self.max_delay_day = MAX_DELAY_DAY  # 最大提取两天预订
        self.open_time = '21:00:00'  # 开放预订时间是21：00点
        # 配置文件中选择的房间
        self.selected_room = None
        for r in get_rooms():
            if r.get('roomNo', None) == self.config['roomNo']:
                self.selected_room = r
                break
        if not self.selected_room:
            raise Exception(
                "目标房间:{},不存在！".format(config['roomNo']))
        session = requests.Session()
        session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36'
        })
        self.__session = session

    def __get(self, url, params={}):
        full_url = 'http://{host}:{port}{url}'.format(
            host=self.host, port=self.port, url=url)
        return self.__session.get(full_url, params=params)

    def __post(self, url, payload={}):
        full_url = 'http://{host}:{port}{url}'.format(
            host=self.host, port=self.port, url=url)
        return self.__session.post(full_url, data=payload)

    def login(self):
        """ 登录操作, 成功返回 True，失败返回 False """
        # post 表单
        payload = {
            'id': self.config['sid'],
            'pwd': self.config['password'],
            'act': 'login'
        }
        r = self.__post(self.login_url, payload)
        r_json = r.json()
        if (r_json.get('msg', None) == 'ok'):
            logging.info('登录成功')
            return True
        else:
            logging.error('登录失败,请检查账号、密码是否配置正确！')
            return False

    def __get_delay_day(self, delay):
        """ 获取今天 + delay 天数 之后的日期。如：2017-01-01 + 2 = 2017-01-03(返回结果) """
        return time.strftime("%Y-%m-%d", time.localtime(time.time() + 3600 * 24 * delay))

    def booking(self):
        """ 自动预约 """
        # 目标日期，一般为后天
        that_day = self.__get_delay_day(self.config['delayDay'])
        start_time = self.config['startTime']
        end_time = self.config['endTime']
        # 预约 post 的表单中开始时间的格式为1800, 表示18：00，去掉小时和分钟的冒号
        start_time = start_time[0:2] + start_time[3:5]
        end_time = end_time[0:2] + end_time[3:5]
        # 其中如果是08:00, 则为800，去掉前面的0
        if start_time[0] == '0':
            start_time = start_time[1:]
        if end_time[0] == '0':
            end_time = end_time[1:]
        # 最终预约的get表单，注意格式细节
        payload = {
            'dev_id': self.selected_room['devId'],
            'lab_id': self.selected_room['labId'],
            'kind_id': self.selected_room['kindId'],
            'type': 'dev',
            'start': that_day + ' ' + self.config['startTime'],
            'end': that_day + ' ' + self.config['endTime'],
            'start_time': start_time,
            'end_time': end_time,
            'act': 'set_resv'
        }
        # 如果预约的是后天，而且当前时间没有超过9点，则等待。
        if (self.config['delayDay'] == self.max_delay_day):
            # 学校预约今天开放时间
            _open_time = time.strftime(
                '%Y-%m-%d', time.localtime()) + ' ' + self.open_time
            open_time = time.mktime(time.strptime(
                _open_time, '%Y-%m-%d %H:%M:%S'))
            # 开放预约时间 - 当前时间 = 需要等待的时间。等待时间一到，发送 get 请求。
            sleep_time = max((open_time - time.time()), 0)
            logging.info('wait {} s'.format(sleep_time))
            # 开始等待
            time.sleep(sleep_time)

        payload['_'] = str(int(time.time()*1000))
        # 发送预订请求
        r = self.__get(self.booking_url, payload)

        # 解析预订结果
        msg = r.json().get('msg', None)
        if (msg and '操作成功' in msg):
            logging.info('自动预订成功。')
            print({
                'start': payload['start'],
                'end': payload['end'],
                'roomName': self.selected_room['devName']
            })
            return True
        else:
            logging.warning('自动预订失败。')
            logging.warning(r.text)
            return False

    def main(self):
        return self.login() and self.booking()


if __name__ == '__main__':
    config = get_config()
    if check_config(config):
        aa = Booking(config)
        aa.main()
