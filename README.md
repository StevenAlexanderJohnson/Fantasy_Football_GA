# Fantasy_Football_GA
Genetic algorithm that collects player data and statistics to generate a fantasy football team.

This was the first project that I wrote in Go. 
I chose to use a genetic algorithm because there is no great algorithm in order to select the highest possible team score. 
There are thousands upon thousands of players and countless combinations of players. 
To be able to test all of these combinations would require until the end of time. 
That is where the genetic algorithm comes into play.

## Use of the Genetic Algorithm

Generally genetic algorithms (GA) are used to simulate "survival of the fittest". 
Who ever is the most fit (or has the highest score) will survive to pass off their genetic code. 
Well in this example, the team who has the highest score will have the most chance to pass on their genetic code. 
In this case, the genetic code is combinations of players on the team. 

This is all a very animalistic definition of the algorithm. 
In this context, it's more akin to people trading players from their team to the next team. 
After trading a player (or group of players) the person would then check again to score their team. 
If their score increased then they are happy and get to continue on. 
If their score decreased too low, then they could be kicked out of the league. 

## How does it work

The algorithm starts by generating teams. The ```PLAYER_STATS_OUTPUT.txt``` file is player stat data collected from an ESPN API using Python. Python was used because it's simply easier to manage the data with nested objects. It was quicker to use Python than to model out all of the nested objects and map them to Go structs. The Python will handle the API data and save it to the text file using a structure that is defined in Go.

From there the collect_player_data function reads the newly formatted data and creates a pool of the players to select from.

Next the GeneticAlgorithmFunctions (GAF) package handles generating teams. The number of teams to generate can be set using the ```--population_size``` flag.

The GAF then scores every team to determine their ```fitness```, then are sorted from most fit to least fit.

The teams are then split into fitness level to allow the most fit to have a higher chance of trading players, in the hope to receive a better player than they lost. This is called the ```crossover```.

After all finished trading, small ```mutations``` are made to a small percentage of the teams in the hope of mutating one player for a better player. This step is only done on a small percentage of the population so that what progress that has been made will not be lost in case a good team is mutated.

This cycle is repeated over and over until the generation cap set by the ```--generation_count``` flag is reached. At that point the highest scoring team will be returned.


## Disclaimer

I am not a sports guy, I don't know all the rules of how this game works. I rely heavily on the points provided from the ESPN API. I know that the team generation is different for every league, so a simply alter to the ```Generate_Team``` function in the ```GAF``` package will be all you would have to do to alter the structure.