FROM scratch
EXPOSE 8080
ADD bin/linux-amd64-simpleregistry /simpleregistry

# Make sure to set 
# -e BINDER_API_KEY=<BLAH> to the token you want to use for auth

ENTRYPOINT ["/simpleregistry"]
