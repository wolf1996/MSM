import pymongo
from pymongo import MongoClient
from datetime import date
import datetime
import numpy as np
import bson

def new_data(sensor_id, date, value, hs):
    result = {
        "sensor_id": sensor_id,
        "date": datetime.datetime.combine(date, datetime.time.min),
        "value": bson.Int64(value),
        "hash": hs
    }
    return result

def main():
    client = MongoClient('localhost', 27017)
    database = client['testdatabase']
    collection = database['sensor_data']
    today = date.today()
    size = 732
    rg = range(0,size)
    c_id = 1
    hs = 'ff62d9d0cc926a7516e408b4ad1a0537'
    vals = [(i + 1)*500 for i in reversed(rg)]
    noises = np.random.randint(-500, 500, size)
    vals += noises
    print(today)
    datalist = []
    for i in rg:
        today -= datetime.timedelta(days=1)
        val = vals[i]
        datalist.append(new_data(c_id, today, val, hs))
    result = collection.insert_many(datalist)
    print(result)


if __name__ == '__main__':
    main()