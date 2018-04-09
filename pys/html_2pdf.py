import pdfkit
import time

path_wkthmltopdf = r'C:\Program Files (x86)\wkhtmltopdf\bin\wkhtmltopdf.exe'
config = pdfkit.configuration(wkhtmltopdf=path_wkthmltopdf)

start = int(round(time.time() * 1000))
#print (millis)
pdfkit.from_url("file:///C:/Users/Dayday/Documents/python%20-%20html%20to%20pdf/test.html", "out.pdf", configuration=config)
end = int(round(time.time() * 1000))
print (end - start)