Building Blockchain
------------------------------------------

This block chain to have transactional information about people who buy books

- We create Book independent struct  [ID, Title, Author, Publish Data, ISBN]
- When the Book is purchased, other struct BookCheckout [User, checkoutDate, IsGenesis]

In Blockchain, the first block is Genesys block

- Block, this can have lot of info, but we will have important, i.e transactional data, which is transaction of books

	
	[Prev hash | Position | Data | TimeStamp| Hash ] <---- [Prev hash | Position | Data | TimeStamp| Hash ] <--- [Prev hash | Position | Data | TimeStamp| Hash ]
	

Position where it lies in entire blockchain



LLD
---------------------------------------------
4 Structs Book, BookCheckout, Block and then Blockchain


  
/CreateBook



													Validate Hash
												  / 
												 /
									Generate_Hash *
												   \
													\
												Validate Block
									Create Block /
									  ^          \ Genesis Block
									  |	
                                      |
New Blockchain---> Write Block---> Add block


Dev
-------------------------------------------------

main.go


Now While creating struct Block we donot have json, its because the two structs i.e Book and BookCheckout we create that using routes using /new, coming from req,res. coming from user
Block is being created by Golang itself, dont need json information


We are not using any DB. overall goal is to mimic the blockchain.
Blockchain itself is a transactional database having peer-to-peer information

hence we have created struct Blockchain


Struct Book--> new Book 


  



