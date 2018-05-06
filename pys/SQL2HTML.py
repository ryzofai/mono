import pandas as pd
import timeit
import sqlite3
from datetime import datetime

def sqlcon2html():
	con = sqlite3.connect('D:\hyperuploader\db\histdb.db')
	c = con.cursor()
	df = pd.read_sql_query("select APPLTYPE, NMBR, TRXDATE, TRXTYPE,TRXCODE, TRNAMT, RUNBAL,FLAG, OPERID, TERMID \
	  from histdb where APPLTYPE = 'IM' order BY ACCTNMBR ASC, TRXDATE DESC, INTSEQNO DESC limit 100;", con)
	#print(df.info)
	df2 = pd.DataFrame(df)
	print(df.dtypes)
	sum_row = df2[['TRNAMT','RUNBAL']].sum()
	df_sum = pd.DataFrame(data=sum_row).T
	df_sum = df_sum.reindex(columns=df.columns)
	#dfsum = df['TRNAMT'].sum()
	#print(df_sum)
	df_final=df2.append(df_sum,ignore_index=True)
	print(df_final.tail())
	print(df_final.dtypes)
	#df.append(pd.Series(df.sum(),name='Total'))
	#df.append(pd.Series(['a', 'b'], index=['col1','col2']), ignore_index=True)
	#print(df)
	#df_final.fillna(0, inplace=True)
	#df_final.fillna("", inplace=True)
	#print(df_final)
	#df_final.to_html('out2.html', index=False, float_format='%.0f')
	df_final.to_html(float_format='%.0f')
	#df2.to_html('out3.html', index=False)
	print(str(df_final).replace('NaN',''))

t = timeit.Timer(writeToHtml_Test)
print('writeToHtml_Test (seconds): ' + str(t.timeit(1)))
