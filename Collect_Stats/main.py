import json
import requests
import threading
import queue

class Stats:
    receivingTouchdowns = 0
    puntReturnTouchdowns = 0
    fumbleTouchdowns = 0
    rushingTouchdowns = 0
    passingTouchdowns = 0
    twoPointRecConvs = 0
    twoPointPassConv = 0
    rushing = 0
    receiving = 0
    passing = 0
    fumblesLost = 0
    interceptions = 0
    extraPoints = 0
    twoPointRushConvs = 0
    fieldGoals = 0
    sacks = 0
    safeties = 0
    defensiveTouchdowns = 0
    blockedKicks = 0

class Player:
    def __init__(self, fullName, weight, height, age, position, stats, years, team_name) -> None:
        self.fullName = fullName
        self.weight = weight
        self.height = height
        self.age = age
        self.years = years
        self.position = position
        self.team = team_name
        self.stats = stats.__dict__


if __name__ == "__main__":
    q = queue.Queue(maxsize=100)

    def get_player_names(player_link):
        player_response = requests.get(player_link)
        player_body = player_response.json()
        if not "statistics" in player_body:
            q.put("SKIP_PLAYER")
            return
        player_statistics = requests.get(player_body["statistics"]["$ref"])
        stats_body = player_statistics.json()
        data_found = False
        player_stat_object = Stats()
        for category in stats_body['splits']['categories']:
            for stat in category['stats']:
                data_found = True
                for stat in category["stats"]:
                    if stat["name"] == "receivingTouchdowns" and stat["value"] != 0:
                        player_stat_object.receivingTouchdowns = stat["value"]
                    elif stat["name"] == "blockedPuntTouchdowns" and stat["value"] != 0:
                        player_stat_object.puntReturnTouchdowns = stat["value"]
                    elif stat["name"] == "fumblesTouchdowns" and stat["value"] != 0:
                        player_stat_object.fumbleTouchdowns = stat["value"]
                    elif stat["name"] == "rushingTouchdowns" and stat["value"] != 0:
                        player_stat_object.rushingTouchdowns = stat["value"]
                    elif stat["name"] == "passingTouchdowns" and stat["value"] != 0:
                        player_stat_object.passingTouchdowns = stat["value"]
                    elif stat["name"] == "twoPointRecConvs" and stat["value"] != 0:
                        player_stat_object.twoPointRecConvs = stat["value"]
                    elif stat["name"] == "twoPointRushConvs" and stat["value"] != 0:
                        player_stat_object.twoPointRushConvs = stat["value"]
                    elif stat["name"] == "rushingYards" and stat["value"] != 0:
                        player_stat_object.rushing = stat["value"]
                    elif stat["name"] == "receivingYards" and stat["value"] != 0:
                        player_stat_object.receiving = stat["value"]
                    elif stat["name"] == "passing" and stat["value"] != 0:
                        player_stat_object.passing = stat["value"]
                    elif stat["name"] == "fumblesLost" and stat["value"] != 0:
                        player_stat_object.fumblesLost = stat["value"]
                    elif stat["name"] == "interceptions" and stat["value"] != 0:
                        player_stat_object.interceptions = stat["value"]
                    elif stat["name"] == "kickExtraPoint" and stat["value"] != 0:
                        player_stat_object.extraPoints = stat["value"]
                    elif stat["name"] == "twoPointPassConvs" and stat["value"] != 0:
                        player_stat_object.twoPointPassConv 
                    elif stat["name"] == "fieldGoals" and stat["value"] != 0:
                        player_stat_object.fieldGoals = stat["value"]
                    elif stat["name"] == "sacks" and stat["value"] != 0:
                        player_stat_object.sacks = stat["value"]
                    elif stat["name"] == "safeties" and stat["value"] != 0:
                        player_stat_object.safeties = stat["value"]
                    elif stat["name"] == "defensiveTouchdowns" and stat["value"] != 0:
                        player_stat_object.defensiveTouchdowns = stat["value"]
                    elif stat["name"] == "kicksBlocked" and stat["value"] != 0:
                        player_stat_object.blockedKicks = stat["value"]
        team_data = requests.get(player_body["team"]["$ref"])
        team_body = team_data.json()
                    
        player_object = Player(player_body["fullName"], player_body["weight"], player_body["height"], player_body["age"], player_body["position"]["name"], player_stat_object, player_body['experience']['years'], team_body["displayName"])
        if not data_found:
            q.put("INVALID STATS")
            return
        q.put(player_object)

    page_index = 1

    with open("./Saved_Data/PLAYER_STATS_OUTPUT.txt", "w+") as outputFile:
        while True:
            response = requests.get("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes?limit=1000&active=true&page=" + str(page_index))
            response_body = response.json()
            if len(response_body["items"]) == 0:
                break
            print("Reading page: " + str(page_index))
            athelete_count = 0
            for athelete in response_body['items']:
                athelete_thread = threading.Thread(target=get_player_names, args=(athelete['$ref'],))
                athelete_thread.start()
                athelete_count += 1
            for i in range(athelete_count):
                thread_value = q.get()
                if(thread_value == "SKIP_PLAYER"):
                    continue
                elif (thread_value == "INVALID STATS"):
                    print("Invalid stats")
                    continue
                outputFile.write(json.dumps(thread_value.__dict__) + "\n")
            thread_pool = []
            page_index += 1
