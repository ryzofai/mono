import win32com.client

word = win32com.client.Dispatch('Word.Application')

doc = word.Documents.Add('C:\\Users\\Dayday\\Documents\\html to doc\\test.html')
doc.SaveAs('C:\\Users\\Dayday\\Documents\\html to doc\\example.doc', FileFormat=0)
doc.Close()

word.Quit()