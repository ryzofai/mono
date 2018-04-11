import twitter

api = twitter.Api(consumer_key='Ni9rZZFOiAHPgnalGxauIReov',
                  consumer_secret='Om3d24bW7NYF3Xst5NxzOXNy51jAV7QUUTB1twYPMWfpLYCU5i',
                  access_token_key='984061142724063233-lvDldadmfua6h7yWHDjNYNb55v1oT5A',
                  access_token_secret='Xz3n1fAw2A5Dn9lapb68L1f2ORL0tussncqT3Eg98oQPK')

try:
    status = api.PostUpdate('Hello from LBP.SRE Twitter app! \nThis app is still in Beta.')
except UnicodeDecodeError:
    print("Your message could not be encoded.  Perhaps it contains non-ASCII characters? ")
    print("Try explicitly specifying the encoding with the --encoding flag")
    sys.exit(2)

    print("{0} just posted: {1}".format(status.user.name, status.text))