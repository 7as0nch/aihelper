import pandas as pd
import os

# 读取CSV文件，不指定表头
df = pd.read_csv('questions.csv', header=None, encoding='utf-8')
print(f"原始数据形状: {df.shape}")
print("前5行:")
print(df.head())

# 提取第一列
first_column = df.iloc[:, 0]
print(f"\n第一列前10个值:")
print(first_column.head(10))

# 保存到新的CSV文件
output_file = 'first_column.csv'
first_column.to_csv(output_file, index=False, header=False, encoding='utf-8')
print(f"\n已保存到 {output_file}")

# 验证新文件
if os.path.exists(output_file):
    with open(output_file, 'r', encoding='utf-8') as f:
        lines = f.readlines()
        print(f"\n新文件行数: {len(lines)}")
        print("新文件前10行:")
        for i, line in enumerate(lines[:10]):
            print(f"{i+1}: {line.strip()}")
else:
    print("输出文件未创建")