# Problem Statement

## Programming Assignment
You will write a program that calculates and provides certain statistics about the market data coming from the Binance crypto exchange.
The requirements of this program are as follows:

- Program shall connect to the raw trade stream WebSocket endpoint: https://binance-docs.github.io/apidocs/spot/en/#trade-streams
- Program shall use as few WebSocket connections as possible in order to include the raw trade stream of every symbol available on Binance. There are currently 1946 symbols provided by Binance, but this changes from time to time. You can fetch the full list of symbols with this REST endpoint: https://binance-docs.github.io/apidocs/spot/en/#exchange-information
- In order to use as few WebSocket connections as possible, you will need to construct a "combined stream", as described here: https://binance-docs.github.io/apidocs/spot/en/#websocket-market-streams. If you are able to make a single combined stream with all of the symbols, your job will be easy, but if not, you will have to determine the max number of symbols you can support per WebSocket connection and manage them all.
- Your program will calculate statistics on a per-symbol basis; that is to say, although you will have data concerning multiple symbols multiplexed onto a single WebSocket connection, your program must be able to calculate the statistics on only the subset of data that apply to any particular symbol.
- Your program should calculate the median of all the prices ("p:" field in the json) seen for a particular symbol up to and including that point. For instance, if we have seen the following prices for the BNBBTC symbol: 3, 4, 2, 1, 5, 8, and then a 2.5 arrives, then the median at that moment would be 3.
- For the avoidance of doubt, a median is calculated as the middle point when all of the data are sorted, in the case of an odd number of points (e.g. the median of 1, 3, 2 is 2.); and it is the mean of the two middle-most data points in the case when there are an even number of points (e.g. the median of 1, 3, 2, 4 is 2.5).

Calculating a median at every new datum is called an "infinite median". The challenge with infinite medians is that the amount of data can grow very large and it will start to get very slow to calculate once the program has seen sufficient data using the naive implementation. There are other methods of calculating the median that are imperfect but much more efficient. Your job is to discover the best of those. To guide your discovery, consider algorithms that would work for many billions of data points on a regular desktop computer.

The time complexity to update your infinite medians must not be worse than O(log n). The time complexity to retrieve the current median for a particular symbol must not be worse than O(1). This should give you some indication of which implementations will not be acceptible.

In order to maintain separate statistics for each symbol, you will need some way of mapping a symbol onto an object in which you will maintain the information you need to calculate the statistics. The natural way to do this is with some form of a hashmap, where symbols are the keys and the values are the objects you need to make the statistics caluclations.

You must implement this hashmap yourself, and the hash function must be perfect. 
That is to say, there can be no collisions across the complete set of symbols supported by your program. For instance, if your hash function takes as input the symbol name and outputs the index in an array, then no other valid symbol can hash to the same index.

## REST API:
You shall make these statistics that you are maintaining accessible via a simple REST API. This API should be public, and doesn't need any form of authentication.

- GET /symbols shall return a json list of all the symbols supported by your API. This should be the same list of symbols as is supported at binance.com.
- GET /&lt;symbol&gt; shall return a json with i) the symbol, ii) the number of data that have been seen for that symbol so far, iii) the infinite median of that symbol's prices, and iv) the most recent price seen for that symbol.

The program should be runnable from the command line, and can be stopped with Ctrl-C.

### Some general reminders:
- Comment your code
- Commit in small chunks and push regularly so we can see how things are going as you work on the assignment.
- Your program does not need to have a UI, but feel free to include one if you find it helpful.
- Be sure to provide very thorough instructions for how to build and run your program. If we cannot run your program, then we will not be able to advance you to the second stage of the interview process.
- You have 24hrs to complete this assignment.
- When you are done please explain in words how your median algorithm works, and where we can find the hashmap implementation.

Ready, set, go!
