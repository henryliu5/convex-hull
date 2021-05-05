import math
import random
import numpy as np
import matplotlib.pyplot as plt
import math
import sys

def generate_circle(num_points, percentage):
    x_list=list()
    y_list=list()
    #Generates hull points
    num_hull = int(num_points * percentage)
    for i in range(num_hull):
        theta = random.random() * 2 * math.pi
        r = 1 + random.random() * 0.05
        x,y = r * math.cos(theta), r * math.sin(theta)
        x_list.append(x)
        y_list.append(y)
    
    num_non_hull = num_points - num_hull
    for i in range(num_non_hull):
        theta = random.random() * 2 * math.pi
        r = random.random() * 1
        x,y = r * math.cos(theta), r * math.sin(theta)
        x_list.append(x)
        y_list.append(y)
    return x_list,y_list

def generate_uniform(num_points):
    x_list=list()
    y_list=list()
    #Generates hull points
    for i in range((num_points)):
        x,y = random.random() * 2 - 1, random.random() * 2 - 1
        x_list.append(x)
        y_list.append(y)
    return x_list,y_list

def generate_normal(num_points):
    x_list=list()
    y_list=list()
    #Generates hull points
    for i in range((num_points)):
        x,y = np.random.normal(), np.random.normal()
        x_list.append(x)
        y_list.append(y)
    return x_list,y_list

def generate_normal_circle(num_points):
    x_list=list()
    y_list=list()
    #Generates hull points
    for i in range((num_points)):
        r,theta = np.random.normal(), random.random() * 2 * math.pi
        x,y = r * math.cos(theta), r * math.sin(theta)
        x_list.append(x)
        y_list.append(y)
    return x_list,y_list

def is_float_try(str):
    try:
        float(str)
        return True
    except ValueError:
        return False

if (len(sys.argv) < 4):
    print("Arguments needed (1) number of points (2) type of points to generate - see README (3) directory to save test case")
num_points = int(sys.argv[1])
point_type = sys.argv[2].upper()
save_dir = sys.argv[3]

if (point_type == "UNI"):
    x,y = generate_uniform(num_points)
elif (point_type == "NOR"):
    x,y = generate_uniform(num_points)
elif (point_type == "NCI"):
    x,y = generate_normal_circle(num_points)
elif (point_type == "LCI"):
    x,y = generate_circle(num_points, 0.1)
elif (point_type == "MCI"):
    x,y = generate_circle(num_points, 0.3)
elif (point_type == "HCI"):
    x,y = generate_circle(num_points, 0.6)
elif (is_float_try(point_type)):
    per = float(point_type)
    x,y = generate_circle(num_points, per)
else:
    x,y = list(), list()

f = open(save_dir, "w")
for i in range(len(x)):
    f.write(str(x[i]) + "," + str(y[i]) + "\n")
f.close()
