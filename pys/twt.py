import twitter

api = twitter.Api(consumer_key='',
                  consumer_secret='',
                  access_token_key='',
                  access_token_secret='')

try:
    status = api.PostUpdate('Hello from SRE Twitter app! \nThis app is still in Beta.')
except UnicodeDecodeError:
    print("Your message could not be encoded.  Perhaps it contains non-ASCII characters? ")
    print("Try explicitly specifying the encoding with the --encoding flag")
    sys.exit(2)

    print("{0} just posted: {1}".format(status.user.name, status.text))
