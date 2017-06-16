import postgresql
from datetime import date
import datetime
import numpy as np


def main():
    db = postgresql.open('pq://iot_api_user:qwerty@localhost:5433/testgodb')
    query = '''
    INSERT INTO DATA VALUES (
        $1,
        $2,
        $3,
        $4
    );
    '''
    today = date.today()
    size = 365
    rg = range(0,size)
    c_id = 1
    hs = 'ff62d9d0cc926a7516e408b4ad1a0537'
    vals = [(i + 1)*500 for i in reversed(rg)]
    noises = np.random.randint(-500, 500, size)
    vals += noises
    print(today)
    p = db.prepare(query)
    for i in rg:
        today -= datetime.timedelta(days=1)
        val = vals[i]
        p(c_id, today, val, hs)
    db.close()
    pass

if __name__ == '__main__':
    main()