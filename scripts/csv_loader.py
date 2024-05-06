from datetime import datetime
import sys
import csv
import json

class Date:
    def __init__(self):
        self.ownerId = None
        self.firstName = None
        self.lastName = None
        self.ethnicity = []
        self.occupation = None
        self.places = []
        self.result = None
        self.age = None
        self.date = None
        self.split = None

class Place:
    def __init__(self):
        self.place = None
        self.cost = None
        self.typeOfPlace = None

class DateEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, Date) or isinstance(obj, Place):
            date = {}
            for k, v in vars(obj).items():
                date[k] = v
            return date
        return super().default(obj)

def strToBool(s):
    if s.lower() in ('true', 't', '1', 'yes'):
        return True
    return False

def read_csv_file(csv_file):
    data = []
    with open('seed.json', 'w') as json_file:
        with open(csv_file, 'r', newline='') as file:
            reader = csv.DictReader(file)
            for row in reader:
                new_place = Place()
                new_place.place = row['place']
                new_place.typeOfPlace = row['typeOfDate']
                if row['cost']:
                    new_place.cost = { "$numberDouble": row['cost'] }
                else:
                    new_place.cost = { "$numberDouble": 0 }
                if row['firstName']:
                    new_date = Date()
                    new_date.ownerId = {"$oid": row['ownerId']}
                    new_date.firstName = row['firstName']
                    new_date.lastName = ""
                    new_date.ethnicity = row['ethnicity'].split('|')
                    new_date.occupation = row['occupation']
                    new_date.result = row['result']
                    new_date.age = int(row['age'])
                    new_date.split = strToBool(row['split'])
                    new_date.date = { "$date": datetime.strptime(row['date'], "%Y-%m-%d").isoformat() }
                    new_date.places.append(new_place)
                    data.append(new_date)
                else:
                    data[-1].places.append(new_place)
        for date in data:
            acc = 0
            for place in date.places:
                acc += float(place.cost["$numberDouble"])
            date.cost = acc
        json.dump(data, json_file, indent=4, cls=DateEncoder)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_csv_file>")
        sys.exit(1)

    csv_file = sys.argv[1]
    read_csv_file(csv_file)