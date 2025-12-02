#!/usr/bin/env python3
"""
CSV文件第一列提取工具
功能：读取CSV文件，提取第一列数据并保存到新文件
"""

import pandas as pd
import os
import sys
import csv
from pathlib import Path

def detect_csv_format(file_path):
    """检测CSV文件的格式（编码、分隔符、是否有表头）"""
    print(f"正在检测文件格式: {file_path}")
    
    # 尝试的编码列表
    encodings = ['utf-8-sig', 'utf-8', 'gbk', 'gb2312', 'latin1']
    
    for encoding in encodings:
        try:
            # 读取文件前几行进行分析
            with open(file_path, 'r', encoding=encoding) as f:
                # 读取前5行
                lines = []
                for i in range(5):
                    line = f.readline()
                    if not line:
                        break
                    lines.append(line)
                
                if not lines:
                    continue
                
                # 尝试检测分隔符
                sample = ''.join(lines)
                sniffer = csv.Sniffer()
                
                try:
                    dialect = sniffer.sniff(sample)
                    delimiter = dialect.delimiter
                    has_header = sniffer.has_header(sample)
                    
                    print(f"✓ 成功检测格式:")
                    print(f"  编码: {encoding}")
                    print(f"  分隔符: {repr(delimiter)}")
                    print(f"  是否有表头: {has_header}")
                    print(f"  样本预览: {lines[0][:50]}...")
                    
                    return {
                        'encoding': encoding,
                        'delimiter': delimiter,
                        'has_header': has_header,
                        'sample_lines': lines
                    }
                except:
                    # 如果自动检测失败，尝试常见分隔符
                    common_delimiters = [',', ';', '\t', '|']
                    for delim in common_delimiters:
                        if delim in lines[0]:
                            print(f"✓ 检测到可能的分隔符 {repr(delim)} (编码: {encoding})")
                            return {
                                'encoding': encoding,
                                'delimiter': delim,
                                'has_header': None,  # 未知
                                'sample_lines': lines
                            }
        except Exception as e:
            continue
    
    print("✗ 无法自动检测文件格式")
    return None

def extract_first_column(input_file, output_file):
    """提取CSV文件的第一列"""
    print(f"\n{'='*60}")
    print("CSV文件第一列提取工具")
    print(f"{'='*60}")
    
    # 检查输入文件是否存在
    if not os.path.exists(input_file):
        print(f"错误: 输入文件不存在: {input_file}")
        return False
    
    # 检测文件格式
    format_info = detect_csv_format(input_file)
    
    if not format_info:
        print("尝试使用默认设置读取文件...")
        format_info = {
            'encoding': 'utf-8-sig',
            'delimiter': ',',
            'has_header': 'infer'  # 让pandas自动推断
        }
    
    try:
        # 读取CSV文件
        print(f"\n正在读取文件: {input_file}")
        
        read_params = {
            'filepath_or_buffer': input_file,
            'encoding': format_info['encoding'],
            'sep': format_info['delimiter']
        }
        
        # 如果有表头信息，使用它
        if format_info.get('has_header') is not None:
            if format_info['has_header']:
                read_params['header'] = 0  # 第一行是表头
            else:
                read_params['header'] = None  # 没有表头
        else:
            read_params['header'] = 'infer'  # 自动推断
        
        data = pd.read_csv(**read_params)
        
        print(f"✓ 成功读取文件")
        print(f"  数据形状: {data.shape} (行数: {data.shape[0]}, 列数: {data.shape[1]})")
        
        # 显示列信息
        if read_params['header'] is None:
            print(f"  列名: 无表头，使用默认列名")
            print(f"  列名列表: {list(data.columns)}")
        else:
            print(f"  列名: {list(data.columns)}")
        
        # 提取第一列
        print(f"\n正在提取第一列数据...")
        first_column_name = data.columns[0]
        first_column_data = data.iloc[:, 0]
        
        print(f"  第一列名称: {first_column_name}")
        print(f"  第一列数据行数: {len(first_column_data)}")
        print(f"  前5个值: {list(first_column_data.head())}")
        
        # 创建新的DataFrame
        new_df = pd.DataFrame({first_column_name: first_column_data})
        
        # 保存到新文件
        print(f"\n正在保存到: {output_file}")
        new_df.to_csv(output_file, index=False, encoding='utf-8')
        
        # 验证保存的文件
        if os.path.exists(output_file):
            file_size = os.path.getsize(output_file)
            print(f"✓ 文件保存成功!")
            print(f"  文件大小: {file_size} 字节")
            print(f"  文件路径: {os.path.abspath(output_file)}")
            
            # 读取保存的文件进行验证
            saved_data = pd.read_csv(output_file, encoding='utf-8')
            print(f"  验证 - 行数: {len(saved_data)}")
            print(f"  验证 - 列名: {list(saved_data.columns)}")
            
            # 显示新文件内容预览
            print(f"\n{'='*60}")
            print("新文件内容预览:")
            print(f"{'='*60}")
            print(saved_data.head(10).to_string(index=False))
            
            if len(saved_data) > 10:
                print(f"... 还有 {len(saved_data) - 10} 行")
            
            return True
        else:
            print("✗ 文件保存失败!")
            return False
            
    except Exception as e:
        print(f"✗ 处理文件时出错: {e}")
        import traceback
        traceback.print_exc()
        return False

def main():
    """主函数"""
    # 文件路径
    input_file = "questions.csv"
    output_file = "questions_first_column.csv"
    
    print(f"输入文件: {input_file}")
    print(f"输出文件: {output_file}")
    
    # 执行提取操作
    success = extract_first_column(input_file, output_file)
    
    if success:
        print(f"\n{'='*60}")
        print("任务完成!")
        print(f"{'='*60}")
        print(f"✓ 已成功提取第一列数据")
        print(f"✓ 新文件: {output_file}")
        print(f"✓ 文件位置: {os.path.abspath(output_file)}")
    else:
        print(f"\n{'='*60}")
        print("任务失败!")
        print(f"{'='*60}")
        sys.exit(1)

if __name__ == "__main__":
    main()
