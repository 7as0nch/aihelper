import csv

with open('questions.csv', 'r', encoding='utf-8-sig') as f:
    reader = csv.reader(f)
    rows = list(reader)
    print('总行数:', len(rows))
    for i in range(min(5, len(rows))):
        print(f'第{i}行: {rows[i]}')
        print(f'字段数: {len(rows[i])}')
        if len(rows[i]) > 0:
            print(f'第一个字段: {rows[i][0]}')
        print()