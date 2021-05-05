import os
import subprocess
import sys
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np


def plot_line(x,y,lab):
    plt.plot(x,y, label=lab)

def plot_point(avg, lab):
    print(avg)
    plt.plot([1,1024], [avg,avg], label=lab)
    
def plot_label(xlab, ylab, title):
    plt.xlabel(xlab)
    plt.ylabel(ylab)
    plt.title(title)
    plt.legend(loc="upper left")

def plot_speedup(base, x, y, lab):
    plt.plot(x,base / y, label=lab)

def plot_speedup_dual(ser_df, para_df, lab):
    ser_times = ser_df.groupby('Var').mean()
    para_times = para_df.groupby('Var').mean()
    speed_up = ser_times["Time"] / para_times["Time"]
    plt.plot(speed_up.index * 100, speed_up.values, label=lab)

trials = 5
#sizes = [100000,200000,400000,800000,1600000,3200000]
#sizes = [100000,200000]
percentage = [0.10, 0.20, 0.30, 0.40, 0.50, 0.60, 0.70, 0.80, 0.90, 1.00]

test_descriptor = str(sys.argv[1])
'''
for p in percentage:
    print("DOing",p)
    s=1600000
    #s=100
    test_gen = "python ../test_case_generation/test_case_gen.py " + str(s) + " " + str(p) + " temp.txt"
    os.system(test_gen)
    print("Genned test")

    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/hull_portion_results/" + test_descriptor + ".txt --voi="+str(p)
    result = subprocess.check_output(cmd, shell=True)
'''
res_df = pd.read_csv("../data_results/hull_portion_results/" + test_descriptor + ".txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
algos = res_df["Algo"].unique()
for algo in algos:
    if ("qh" not in algo and "chan" not in algo):
        continue
    algo_df = res_df[res_df["Algo"] == algo]
    plot_line(algo_df["Var"] * 100, algo_df["Time"] / 1000000, algo)

plot_label("Percentage of Circumference Points", "Time (ms)", "Convex Hull Algorithm Time Hull Vary")
plt.savefig("../data_results/hull_portion_results/" + test_descriptor + ".png")
plt.close()

for algo in algos:
    if ("qh" not in algo and "chan" not in algo):
        continue
    para_algo = algo.replace("serial", "parallel") 
    if ("serial" in algo and para_algo in algos):
        serial_df = res_df[res_df["Algo"] == "serial_qh"]
        para_df = res_df[res_df["Algo"] == para_algo]
        plot_speedup_dual(serial_df, para_df, algo.replace("_serial", "").strip())

plot_label("Percentage of Circumference Points", "Speedup", "Convex Hull Algorithm Speedup Hull Vary")
plt.savefig("../data_results/hull_portion_results/" + test_descriptor + "_speedup.png")