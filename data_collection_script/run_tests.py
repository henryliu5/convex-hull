import os

os.system("python size_time_collect.py LCI LowHull")
os.system("python size_time_collect.py MCI MedHull")
os.system("python size_time_collect.py HCI HighHull")
os.system("python hull_perc_collect.py HullVary")