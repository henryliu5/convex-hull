import os

os.system("python size_time_collect.py UNI Uniform")
os.system("python size_time_collect.py LCI LowHull")
os.system("python size_time_collect.py MCI MedHull")
os.system("python size_time_collect.py HCI HighHull")
os.system("python hull_perc_collect.py HullVary")
os.system("python num_procs_collect.py Procs")
os.system("python coalesce_chans.py coalesce")