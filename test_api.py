# coding:utf-8
import sys
import urllib3, json, base64, time, hashlib
from datetime import datetime

urllib3.disable_warnings()

with open("doc/exchainge.png", 'rb') as f:
    img_data = f.read()
img_data = base64.b64encode(img_data).decode('utf-8')

# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys())
    if 'data' in name_list: # data 按 key 排序, 中文不进行性转义，与go保持一致
        param['data'] = json.dumps(param['data'], sort_keys=True, ensure_ascii=False).replace(' ','')
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])


if __name__ == '__main__':
    if len(sys.argv)<2:
        print("usage: python3 %s <host> <port>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    port = sys.argv[2]

    body = {
        'version'  : '1',
        'sign_type' : 'SHA256', 
        'data'     : {
            'userkey'   : 'poCX4ig37Ljq4nTX6DjQlz9EgITtIhAONaok8PJ7fzw=',
            'userkey_a' : 'poCX4ig37Ljq4nTX6DjQlz9EgITtIhAONaok8PJ7fzw=',
            'userkey_b' : 'c4l1SkZiwkQsULUah6OPuei0LabubLxnni9tLM1T3Tk=',
            'assets_id' : '123',
            'data'      : 'zzzzzxxxxxxx', # img_data
            'user_name' : '测试2',
            'user_type' : 'buyer',
            'block_id'  : '9681ddbe-3830-44d1-b8af-a41f73a1346a', # 85c2d455-755d-461b-89fc-8a9327f8223a
        }
    }

    secret = 'MjdjNGQxNGU3NjA1OWI0MGVmODIyN2FkOTEwYTViNDQzYTNjNTIyNSAgLQo='
    appid = '4fcf3871f4a023712bec9ed44ee4b709'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appid'] = appid

    param_str = gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, secret)

    if body['sign_type'] == 'SHA256':
        sha256 = hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')
        signature_str =  base64.b64encode(sha256).decode('utf-8')
    else: # SM2
        #signature_str = sm2.SM2withSM3_sign_base64(sign_str)
        pass

    #print(sign_str.encode('utf-8'))
    #print(sha256)
    #print(signature_str)

    body['sign_data'] = signature_str

    body = json.dumps(body)
    #print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    host = 'http://%s:%s'%(hostname, port)
    #url = host+'/api/query_by_assets'
    #url = host+'/api/biz_register'
    #url = host+'/api/biz_contract'
    url = host+'/api/query_block'

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
