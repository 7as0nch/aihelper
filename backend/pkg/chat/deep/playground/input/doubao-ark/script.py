import pandas as pd

df = pd.read_csv('questions.csv')
first_column = df.iloc[:, 0]
first_column.to_csv('output.csv', index=False)