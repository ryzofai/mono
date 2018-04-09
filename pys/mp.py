from multiprocessing import Process, Queue
import time
import random

def job(q):
	print('job 1 started')
	while True:
		if q.empty() is False:
			print('job1 ' + q.get())
			time.sleep(1)
	
def job2(q):
	print('job 2 started')
	while True:
		if q.empty() is False:
			print('job2 ' + q.get())
			time.sleep(2)
	
def jobQ(q):
	print('job Q started')
	#while True:
	for i in range (1, 20):
		q.put('message: ' + str(random.randint(1,1000000)))
		time.sleep(0.5)
		#print(q.get())
	
	time.sleep(10)
	print('second wave waking up')
	
	for i in range (1, 20):
		q.put('message w2: ' + str(random.randint(1,1000000)))
		time.sleep(0.5)
		
if __name__ == '__main__':
	q = Queue()
	p1 = Process(target=job, args=(q,))
	p2 = Process(target=job2, args=(q,))
	p3 = Process(target=jobQ, args=(q,))
	
	p3.start()
	p1.start()
	p2.start()
	

		