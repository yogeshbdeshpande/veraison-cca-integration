~~~aasvg
                 .----------.
                |  Artifact  |
                 '----+-----'
                      v
                 .----+----.  .----------.  Identifiers
Issuer      --> | Statement ||  Envelope  +<------------------.
                 '----+----'  '-----+----'                     |
                      |             |           +--------------+---+
                       '----. .----'            | Identity         |
                             |                  | Documents        |-------
                             v                  +-------+------+---+      |
                        .----+----.                     |                 |
                       |  Signed   |    COSE Signing    |                 |
                       | Statement +<-------------------+                 |
                        '----+----'            +-----------+----+         |    
                             |               +-+---------+----+ |         |
                          .-' '------------->+ Transparency   | |         |
                         |    .----------.   | Services       | |         |
Transparency -->         |  .-+--------. |   |                |-+         |
     Service             |  | Receipts +<----+----+------+----+           |
                         |  '-+--+-----'                  |               |
						 |                                |               |
                          '-. .-'                         |               |
                             |                            |               |
                             v                            |               |
                       .-----+-----.                      |               |
                      | Transparent |                     |               |
                      |  Statement  |                     |               |
                       '-----+-----'                      |               |
                             |                            |               |
                             |'-------.     .-------------)---------------'
                             |         |   |              |
                             |         v   v              |  
                             |    .----+---+-----------.  |
Verifier     -->             |   / Verify Transparent /   |
                             |  /      Statement     /    |
                             | '--------------------'     |
                             v                            v
                    .--------+---------.      .-----------+-----.
Auditor      -->   / Collect Receipts /      /   Replay Log    /
                  '------------------'      '-----------------'
~~~
