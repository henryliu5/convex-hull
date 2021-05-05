import os
import subprocess
import sys
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
    plt.plot(speed_up.index, speed_up.values, label=lab)

trials = 5
#sizes = [100000,200000,400000,800000,1600000,3200000]
#sizes = [100000,200000]

test_descriptor = str(sys.argv[1])
l = [False, True]
for p in l:
    print("Trying",l)
    s=1600
    #s=100
    test_gen = "python ../test_case_generation/test_case_gen.py " + str(s) + " " + "HCI" + " temp.txt"
    os.system(test_gen)

    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/coalesce_folder/" + test_descriptor + ".txt --voi="+str(p) + " --coalesce=" + str(p)
    result = subprocess.check_output(cmd, shell=True)

res_df = pd.read_csv("../data_results/coalesce_folder/" + test_descriptor + ".txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
res_df = res_df[(res_df["Algo"]=="serial_chans") | (res_df["Algo"]=="parallel_chans")]
print(res_df)
labels = list()
times = list()
for index, row in res_df.iterrows():
    times.append(row["Time"])
    name = row["Algo"]
    if (row["Var"] == True):
        name = name + " Coalesced"
    labels.append(name)

plt.bar(labels, times, color ='maroon',
        width = 0.4)

plot_label("Algo Ran", "Time (ms)", "Convex Hull Algorithm on 1,600,000 High Circle Sample")
plt.savefig("../data_results/coalesce_folder/" + test_descriptor + ".png")
plt.close()
