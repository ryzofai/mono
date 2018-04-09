import time

# Open a file
fo = open("test.txt", "r")

ctr = 1

while True:
	#str = fo.read(10);
	str = fo.readline()
	if str != "" and str != "\n":
		print ("Read String is : ", str)
	# Close opend file
	#++ctr
	time.sleep(1)
fo.close()