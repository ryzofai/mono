import time

# Open a file
fo = open("test.txt", "r")

while True:
	str = fo.readline()
	if str != "" and str != "\n":
		print ("Read String is : ", str)
	time.sleep(1)
fo.close()