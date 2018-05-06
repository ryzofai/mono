import pandas as pd
import timeit
import sqlite3
from datetime import datetime

def convertoCSV():
	fr = open('HISTORYTBL2.DAT', 'r')
	fw = open('history_csv.txt', 'w')

	hline = 'APPLTYPE' + ',' + 'ANR'+ ',' + 'TRXDATE' + ',' + 'TRXTYPE' + ',' + 'TRXCODE' + ',' + \
			'TRNAMT' + ',' + 'RUNBAL' + ',' + 'FLAG' + ',' + 'OPERID' + ',' + 'TERMID' + ',' + 'BRANCHNO' + ',' + 'SEQNUM' + ',' + \
			'POSTDATE' + ',' + 'POSTTIME' + ',' + 'SERNUM' + ',' + 'SOURCE' + ',' + 'STMTSYMBOL' + ',' + 'DSTAMP' + ',' + 'INTSEQNO'

	fw.write(hline + '\n')
	INTSEQNO = 1
	for line in fr:
		
		line.rstrip()
		APPLTYPE = line[0:2]
		ANR = line[2:12]
		TRXDATE = line[12:18]
		TRXTYPE = line[18:20]
		TRXCODE = line[20:24]
		TRNAMT = line[24:39]
		RUNBAL = line[39:54]
		FLAG = line[54:55]

		OPERID = line[55:59]
		TERMID = line[59:63]
		BRANCHNO = line[63:66]
		SEQNUM = line[66:71]
		POSTDATE = line[71:77]
		POSTTIME = line[77:83]
		SERNUM = line[83:93]
		SOURCE = line[93:95]
		STMTSYMBOL = line[95:97]

		wline = APPLTYPE + ',' + ANR + ',' + TRXDATE + ',' + TRXTYPE + ',' + TRXCODE + ',' + \
				TRNAMT + ',' + RUNBAL + ',' + FLAG + ',' + OPERID + ',' + TERMID + ',' + BRANCHNO  + ',' + SEQNUM + ',' + \
				POSTDATE + ',' + POSTTIME + ',' + SERNUM + ',' + SOURCE + ',' + STMTSYMBOL + ',' + str(datetime.now()) + ',' + str(INTSEQNO)
		
		fw.write(wline + '\n')
		INTSEQNO+=1

def uloadToSQLite3():
	con = sqlite3.connect('D:\hyperuploader\db\histdb_testupload.db')
	c = con.cursor()
	file = r'history_csv.txt'
	df = pd.read_csv(file)
	df.to_sql('histdb_testupload', con, if_exists='replace')

def convertToCTL():
	con = sqlite3.connect('D:\hyperuploader\db\histdb.db')
	# sqlite3 D:/hyperuploader/db/histdb.db "select * from histdb;" > hist.csv
	c = con.cursor()
	fw = open('out.html', 'w')
	#file = r'history_csv.txt'
	#df = pd.read_csv(file)
	#df.info()
	#df.to_sql('histdb', con, if_exists='replace')
	#print(pd.read_sql_query("select * from histdb limit 5;", con))
	#print(pd.read_sql_query("select * from histdb where trxcode = '6018' limit 5;", con))
	#print(pd.read_sql_query("select * from histdb where APPLTYPE = 'IM' order BY ANR ASC, TRXDATE DESC, INTSEQNO DESC limit 100;", con))
	#print(pd.read_sql_query("select count(distinct(ANR)) from (select * from histdb where APPLTYPE = 'IM' order BY ANR ASC);", con))
	#df = pd.read_sql_query("select count(distinct(ANR)) from (select * from histdb where APPLTYPE = 'IM' order BY ANR ASC);", con)
	df = pd.read_sql_query("select APPLTYPE, ANR, TRXDATE, TRXTYPE,TRXCODE, TRNAMT, RUNBAL,FLAG, OPERID, TERMID \
	  from histdb where APPLTYPE = 'IM' order BY ANR ASC, TRXDATE DESC, INTSEQNO DESC limit 100;", con)
	#df.to_html('out2.html', index=False)
	for index, row in df.iterrows():
		print(row['OPERID'])
		if row['TRXCODE'] != '6018' and  row['TRXCODE'] !='6118':
			print(row)
	#fw.write(HTML(df.to_html(index=False)))
	#print(pd.read_sql_query("select count(*) from histdb;", con))
	#print(pd.read_sql_query("select * from histdb where ANR = 0027201849;", con))
	#df.loc([df['APPLTYPE'] == 'IM'])
	#df.sort_values(by='ANR', ascending=False, by='TRXDATE', ascending=False,)


def writeToHtml_Test():
	con = sqlite3.connect('D:\hyperuploader\db\histdb.db')
	c = con.cursor()
	df = pd.read_sql_query("select APPLTYPE, ANR, TRXDATE, TRXTYPE,TRXCODE, TRNAMT, RUNBAL,FLAG, OPERID, TERMID \
	  from histdb where APPLTYPE = 'IM' order BY ANR ASC, TRXDATE DESC, INTSEQNO DESC limit 100;", con)
	#print(df.info)
	df2 = pd.DataFrame(df)
	print(df.dtypes)
	sum_tranamt = df2[['TRNAMT']].sum()
	sum_runbal = df2[['RUNBAL']].sum()
	print(sum_tranamt.get_value('TRNAMT'))
	print(sum_runbal.get_value)
	#blank_row = pd.DataFrame([''],ignore_index=True)
	#df_sum = pd.DataFrame(data=sum_row).T
	#print(df_sum)
	#new_index= ['aaaaaaaaa', '2 ', '3 ', 'TRNAMT','RUNBAL', '6','7','8','9','1-']
	#df_sum = df_sum.reindex(columns=new_index)
	#dfsum = df['TRNAMT'].sum()
	#print(df_sum)
	#df_final=df2.append(blank_row,ignore_index=False)
	#df_final=df2.append(pd.DataFrame(['']),ignore_index=True)
	#to_html = str(df2.to_html(index=False, float_format='%.0f')).replace('NaN','')
	#to_html = to_html + '<br>'


	#df_final=df2.append(df_sum,ignore_index=True)
	#print(df_final.tail())
	#print(df_final.dtypes)
	#df.append(pd.Series(df.sum(),name='Total'))
	#df.append(pd.Series(['a', 'b'], index=['col1','col2']), ignore_index=True)
	#print(df)
	#df_final.fillna(0, inplace=True)
	#df_final.fillna("", inplace=True)
	#print(df_final)
	#df_final.to_html('out2.html', index=False, float_format='%.0f')
	#df_final.to_html(float_format='%.0f')
	#df2.to_html('out3.html', index=False)
	
	#print(str(df_final).replace('NaN',''))
	finout = str(df2.to_html(index=False, float_format='%.0f')).replace('NaN','')
	final = finout +  '<td>     </td>' + '<td>IM</td>' +str(sum_tranamt.get_value('TRNAMT')) + '   asd' + str(sum_runbal.get_value('RUNBAL')) 
	fw = open('out-fin.html', 'w')
	fw.write(final)

def rwtestFile():
	fr = open('history_csv.txt', 'r')
	fw = open('out.html', 'w')

t = timeit.Timer(writeToHtml_Test)
print('writeToHtml_Test (seconds): ' + str(t.timeit(1)))

#t = timeit.Timer(convertoCSV)  # 219 seconds / 3 mins 39 seconds
#print('convert to csv (seconds): ' + str(t.timeit(1)))

#t2 = timeit.Timer(uloadToSQLite3)  # 245 seconds / 4 minutes
#t = timeit.Timer(pandasDF)
#print('upload to sqlite3 (seconds): ' + str(t2.timeit(1)))

#t = timeit.Timer(rwtestFile)
#print('testfile seconds: ' + str(t.timeit(1)))

# df.to_csv('out.txt',index=False,header=True)

################################################
# FLAG position (55) char,
# OPERID position (56:59) char,
# TERMID	position (60:63) char,
# BRANCHNO position (64:66) char,
# SEQNUM	position (67:71) char,
# POSTDATE position (72:77) date 'MM-DD-YY',
# POSTTIME position (78:83) char,
# SERNUM  position (84:93) char,
# SOURCE	position (94:95) char,
# STMTSYMBOL position (96:97) char,
# DSTAMP sysdate,
# INTSEQNO integer "seqid.nextval"
