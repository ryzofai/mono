import time

# Open a file
f = open("filesplittest.txt", "r")

tstart = time.time()
print(sum(1 for line in open('filesplittest.txt')))
print(time.time() - tstart)

tstart = time.time()

ctr = 0
for line in f:
	ctr+=1

print(ctr)
print(time.time() - tstart)

#print(len(f.read()))

#str = fo.readline()
f.close()