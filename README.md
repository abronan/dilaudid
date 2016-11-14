<pre>

    <b>~*~ hlian/dilaudid ~*~</b>

    encoding and decoding <a href="https://github.com/alizain/ulid/blob/master/README.md">ULID</a>s in go
    
      | everybody is identical in their
      | secret unspoken belief that way
      | deep down they are different
      | from everyone else –DFW
    
    <b>about</b>
    
    • this is a fork of <a href="https://github.com/imdario/go-ulid">imdario/go-ulid</a>

    • we represent ULIDs differently
      (not as bytes)

    • for some reason, ULID libraries
      never include decoding
      • why is that?
      • is it our busy millenial lives?

    • anyway this library decodes
    
    <b>installation</b>
    
    $ go get github.com/hlian/dilaudid
    
    <b>usage</b>
    
    import "github.com/hlian/dilaudid"
    u := dilaudid.NewRandom()
    v, err := dilaudid.Decode("0056190ba946ee68a1c12c3c77b24399")
    
    <b>performance</b>
    
    $ go test --bench

    but also, who cares

</pre>
