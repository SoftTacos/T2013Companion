let socket = new WebSocket("ws://192.168.0.12:8082/dm/socket");
socket.binaryType = 'arraybuffer';
var characterCheckOption = document.getElementById("characterCheckOption")
var statCheckOption = document.getElementById("statCheckOption")
var skillCheckOption = document.getElementById("skillCheckOption")
var bonusCheckOption = document.getElementById("bonusCheckOption")
var difficultyCheckOption = document.getElementById("difficultyCheckOption")


//case switch instead of map because js is spaghetti
function route(response) {
  switch (response.type) {
    case 0:

      break;
    case 10:
      receiveOldChatMessages(response)
      break;
    case 11:
      receiveChatMessage(response)
      break;
    default:
      console.log("ERROR: response type not found. ", response)
      break;
  }
}

function requestStatCheck() {
  //sample skill check
  charName = characterCheckOption.options[characterCheckOption.selectedIndex].value;
  statName = statCheckOption.options[statCheckOption.selectedIndex].text;
  bonus = bonusCheckOption.options[bonusCheckOption.selectedIndex].text;
  difficulty = difficultyCheckOption.options[difficultyCheckOption.selectedIndex].text;

  message = charName + "," + statName + "," + bonus + "," + difficulty
  encReq = encodeRequest(0, message)
  if (encReq != null) {
    socket.send(encReq)
    console.log(message)
  }
}

function requestSkillCheck() {
  //sample skill check
  charName = characterCheckOption.options[characterCheckOption.selectedIndex].value;
  statName = statCheckOption.options[statCheckOption.selectedIndex].text;
  skillName = skillCheckOption.options[skillCheckOption.selectedIndex].text;
  bonus = bonusCheckOption.options[bonusCheckOption.selectedIndex].text;
  difficulty = difficultyCheckOption.options[difficultyCheckOption.selectedIndex].text;

  message = charName + "," + statName + "," + skillName + "," + bonus + "," + difficulty
  encReq = encodeRequest(1, message)
  if (encReq != null) {
    socket.send(encReq)
    console.log(message)
  }
}

