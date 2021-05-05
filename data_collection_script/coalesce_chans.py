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
    plt.plot(speed_up.index, speed_up.values, label=lab)

trials = 5
#sizes = [100000,200000,400000,800000,1600000,3200000]
#sizes = [100000,200000]

test_descriptor = str(sys.argv[1])
for p in range(0,7):
    n = 2**p
    print("DOing",p)
    s=1600000
    test_gen = "python ../test_case_generation/test_case_gen.py " + str(s) + " " + "UNI" + " temp.txt"
    os.system(test_gen)

    print("t2") 
    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/coalesce_folder/" + test_descriptor + "_coal.txt --voi="+str(n) + " --procs=" + str(n) + " --impl=chan --coalesce"
    result = subprocess.check_output(cmd, shell=True)

    print("t1")
    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/coalesce_folder/" + test_descriptor + ".txt --voi="+str(n) + " --procs=" + str(n) + " --impl=chan"
    result = subprocess.check_output(cmd, shell=True)

    print("t3")
    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/coalesce_folder/" + test_descriptor + "_sing_iter.txt --voi="+str(n) + " --procs=" + str(n) + " --impl=chan --simul_iters=1"
    result = subprocess.check_output(cmd, shell=True)

    print("t4")
    cmd = "../runner --trials=" + str(trials) + " --input=./temp.txt --do_output=false --result_file=../data_results/coalesce_folder/" + test_descriptor + "_qh.txt --voi="+str(n) + " --procs=" + str(n) + " --impl=quic"
    result = subprocess.check_output(cmd, shell=True)

res_df = pd.read_csv("../data_results/coalesce_folder/" + test_descriptor + ".txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
print(res_df)
res_df = res_df[res_df["Algo"] == "parallel_chans"]

res_df_co = pd.read_csv("../data_results/coalesce_folder/" + test_descriptor + "_coal.txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
res_df_co = res_df_co[res_df_co["Algo"] == "parallel_chans"]

res_df_si = pd.read_csv("../data_results/coalesce_folder/" + test_descriptor + "_sing_iter.txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
res_df_si = res_df_si[res_df_si["Algo"] == "parallel_chans"]

res_df_qh = pd.read_csv("../data_results/coalesce_folder/" + test_descriptor + "_qh.txt", sep='\t', names=["Algo", "Trials", "Time", "Var"])
best_serial_time = res_df_qh[res_df_qh["Algo"] == "serial_qh"]["Time"].mean()
res_df_qh = res_df_qh[res_df_qh["Algo"] == "parallel_qh"]


plot_line(res_df["Var"], res_df["Time"] / 1000000, "Chan's")
plot_line(res_df_co["Var"], res_df_co["Time"] / 1000000, "Chan's (Coalesce)")
plot_line(res_df_si["Var"], res_df_si["Time"] / 1000000, "Chan's 1 Simul Iter")
plot_line(res_df_qh["Var"], res_df_qh["Time"] / 1000000, "Quickhull")
plot_label("Max Procs", "Time (ms)", "Convex Hull Algorithm Time")
plt.savefig("../data_results/coalesce_folder/" + test_descriptor + ".png")
plt.close()



print(best_serial_time, "best")
plot_line(res_df["Var"], best_serial_time / res_df["Time"] / best_serial_time, "Chan's")
plot_line(res_df_co["Var"], best_serial_time / res_df_co["Time"], "Chan's (Coalesce)")
plot_line(res_df_si["Var"], best_serial_time / res_df_si["Time"], "Chan's 1 Simul Iter")
plot_line(res_df_qh["Var"], best_serial_time / res_df_qh["Time"], "Quickhull")
plot_label("Max Procs", "Speedup", "Convex Hull Algorithm Speedup")
plt.savefig("../data_results/coalesce_folder/" + test_descriptor + "_sp.png")
plt.close()
