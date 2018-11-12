import time
import matplotlib.pyplot as plt

# reddit comments => submission time nicht relevant (comments) https://medium.com/hacking-and-gonzo/how-reddit-ranking-algorithms-work-ef111e33d0d9
# hot sort or Hacker News’s ranking algorithm ==> submission time relevant (haupteintraege) https://medium.com/hacking-and-gonzo/how-hacker-news-ranking-algorithm-works-1d9b0cf2c08d

def calc_score(P,T,G=float(1.8)):
    score = (P-1) / pow((T+2),G)
    return score


def generate_data(P,T,G,t,r):
    x = list()
    y = list()
    for i in range(0,r):
        T += t
        score = calc_score(P,T,G)
        x.append(score)
        y.append(T)
    return y,x

x,y = generate_data(50,0,1.8,0.2,100)
x2,y2 = generate_data(50,0,2.8,0.2,100)
x3,y3 = generate_data(50,0,0.8,0.2,100)



plt.plot(x,y, color='red')
plt.plot(x2,y2, color='blue')
plt.plot(x3,y3, color='orange')
plt.xlabel('score')
plt.ylabel('time(h)')
plt.title('ranking algorythmusB')

plt.show()
