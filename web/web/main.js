// Convert a hex string to a byte array
function hexToBytes(hex) {
    for (var bytes = [], c = 0; c < hex.length; c += 2)
    bytes.push(parseInt(hex.substr(c, 2), 16));
    return bytes;
}

// Convert a byte array to a hex string
function bytesToHex(bytes) {
    for (var hex = [], i = 0; i < bytes.length; i++) {
        hex.push((bytes[i] >>> 4).toString(16));
        hex.push((bytes[i] & 0xF).toString(16));
    }
    return hex.join("");
}

function reverseBytes(txid){
  return bytesToHex(hexToBytes(txid).reverse())
}

function send(){ //217 chars
  var pkey = document.getElementById('pkey').value
  var title = document.getElementById('title').value
  var type = document.getElementById('data_type').value
  var hash = document.getElementById('hash').value
  var prefix = "0xe901" //theca init
  var payload = [hash,type,title]
  console.log(payload);
  var tx = {
      data: [prefix, hash,type,title],
      cash: { key: pkey }
    }
  datacash.send(tx, function(err, res) {
    console.log(res)
  })
}


function reply(txid,comment){ //comment 184chars
  var pkey = document.getElementById('pkey').value
  var commentCounter = document.getElementById("comment_counter_"+txid)
  var commentImg = document.getElementById("comment_"+txid)
  var prefix = "0x6d03" // memo reply
  var reverseByteTxId = "0x"+reverseBytes(txid)
  var tx = {
      data: [prefix, reverseByteTxId,comment],
      cash: { key: pkey }
    }
  datacash.send(tx, function(err, res) {
    if(err != null){
      return false
    }else{
      commentCounter.innerHTML = parseInt(commentCounter.innerHTML) + 1
      commentImg.src = "icons/comment_1.png"
      console.log(res)
      return true
    }
  })
}

function like(txid,counter){
  var pkey = document.getElementById('pkey').value
  var likeCounter = document.getElementById("like_counter_"+txid)
  var likeImg = document.getElementById("like_"+txid)
  var prefix = "0x6d04" //memo like
  var reverseByteTxId = "0x"+reverseBytes(txid)
  var tx = {
      data: [prefix, reverseByteTxId],
      cash: { key: pkey }
    }
  datacash.send(tx, function(err, res) {
    if(err != null){
      return false
    }else{
      // likeImg.className = "liked"
      likeCounter.innerHTML = parseInt(likeCounter.innerHTML) + 1
      // change like img
      likeImg.src = "icons/heart_1.png"
      console.log(res)
      return true
    }
  })
}

function follow(address){
  var pkey = document.getElementById('pkey').value
  var prefix = "0x6d05"
  //todo: validate txid pattern
  var tx = {
      data: [prefix, txid],
      cash: { key: pkey }
    }
  datacash.send(tx, function(err, res) {
    if(err != null){
      return false
    }else{
      return true
    }
  })
}

function unfollow(address){
  var pkey = document.getElementById('pkey').value
  var prefix = "0x6d06"
  //todo: validate txid pattern
  var tx = {
      data: [prefix, txid],
      cash: { key: pkey }
    }
  datacash.send(tx, function(err, res) {
    if(err != null){
      return false
    }else{
      return true
    }
  })
}

function check_link(link){
  return link
}

function check_type(type){
  return type
}

function check_title(title){
  return title
}

function play(hash,title,sender){
  console.log(hash,title);
  download_torrent(hash,title,sender);
}

function getFromAPI() {
  // var search_string = document.getElementById('search').value

  var url = "http://192.168.12.5:8000/api/tx/positions";
  console.log(url)
  var header = {
    headers: { key: "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44" }
  };

  // fetch(url, header).then(function(r) {
  fetch(url).then(function(r) {
    return r.json()
  }).then(function(r) {

    console.log(r)
    document.getElementById('bitdb_output').innerHTML = ""
    document.getElementById('bitdb_output_container').style.display = "block"

    if(r.length != 0){
      var tr = document.createElement('tr');
      for(i in r){
        var tx = r[i]
        list_tx_results(tx,true);
      };
    };
  })
};


function bitdb_get_magnetlinks(limit) {
  var search_string = document.getElementById('search').value


  var query = {
  	"v": 3,
  	"e": { "out.b1": "hex"  },
  	"q": {
  		"db": ["c","u"],
  		"find": {
        "out.s4": { "$regex": search_string, "$options": "i" },
  			"out.b1": "e901",
  			"out.b0": {
  				"op": 106
  			}
  		},
  		"limit":100000,
  		"project": {
  		  "out.b0": 1,
  			"out.b1": 1,
  			"out.s2": 1,
  			"out.s3": 1,
  			"out.s4": 1,
  			"tx": 1,
  			"blk": 1,
  			"in.e.a":1,
  			"_id": 1
  		}
  	}
  };
  var b64 = btoa(JSON.stringify(query));
  var url = "https://bitdb.network/q/" + b64;

  var header = {
    headers: { key: "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44" }
  };

  fetch(url, header).then(function(r) {
    return r.json()
  }).then(function(r) {

    document.getElementById('bitdb_output').innerHTML = ""
    document.getElementById('bitdb_output_container').style.display = "block"

    if(r['c'].length != 0){
      var tr = document.createElement('tr');
      for(i in r['c']){
        var tx = r['c'][i]
        list_tx_results(tx,true);
      };
    };

    if(r['u'].length != 0){
      var tr = document.createElement('tr');
      for(i in r['u']){
        var tx = r['u'][i]
        list_tx_results(tx,false);
      };
    };
  })
};


function openComments(txid){
  console.log(txid);
  var request = new XMLHttpRequest();

  request.open('GET', 'http://192.168.12.5:8000/api/comments/'+txid, true);
  request.onload = function () {
    var modal = document.getElementById("commentsModal");
    var span = document.getElementsByClassName("close")[0];
    var content = document.getElementsByClassName("modal-content")[0];
    var contentReset = "<span class='close'>&times;</span>";
    var title = document.createElement('h1');

    modal.style.display = "block";
    span.onclick = function() {
        modal.style.display = "none";
        content.innerHTML = contentReset;
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
            content.innerHTML = contentReset;
        }
    }

    var data = JSON.parse(this.response);
    console.log(data)
    if (request.status >= 200 && request.status < 400) {
      var ul = document.createElement('ul');
      if(data){
        data.forEach(tx => {
          var li = document.createElement('li');
          li.innerHTML =  "<a class='result-tx-link' data-toggle='tooltip' title='Tx-Data: " + JSON.stringify(tx) + "' target='_blank' href='https://blockchair.com/bitcoin-cash/transaction/"+ tx.Txid +"'><span class='glyphicon glyphicon-th'></span></a> "+tx.Message;
          ul.appendChild(li)
        });
        title.innerHTML = "Comments"
      } else {
        title.innerHTML = "No Comments"
      }

      content.appendChild(title)
      content.appendChild(ul)
    } else {
      console.log('error');
    }
  }
  request.send();
}

function list_tx_results(tx,confirmed){
  var tr = document.createElement('tr');
  var td_txid = document.createElement('td');
  var td_like = document.createElement('td');
  var td_comments = document.createElement('td');
  var td_6a_magnethash = document.createElement('td');
  var td_6a_title = document.createElement('td');
  var td_6a_type = document.createElement('td');
  var td_sender = document.createElement('td');
  var td_blockheight = document.createElement('td');
  var td_play = document.createElement('td');
  var td_score = document.createElement('td');


  td_txid.innerHTML = "<a class='result-tx-link' data-toggle='tooltip' title='Tx-Data: " + JSON.stringify(tx) + "' target='_blank' href='https://blockchair.com/bitcoin-cash/transaction/"+ tx.Txid +"'><span class='glyphicon glyphicon-th'></span></a>";
  td_txid.style.width = "15px";
  if (tx.Likes > 0){
    likeImage = "heart_1.png"
  }else{
    likeImage = "heart_0.png"
  }

  var testcomment = "test comment, bla bla bla!\n dies ist ein test :)"
  if (tx.Comments > 0){
    commentImage = "comment_1.png"
    // td_comments.innerHTML = "<a title='comment' onclick='reply(`"+ tx.Txid +"`,`"+ testcomment +"`)'><img class='comment' id='comment_"+ tx.Txid +"' height='20' src='icons/"+ commentImage +"'><span class='commentcounter' id='comment_counter_"+ tx.Txid +"'>"+ tx.Comments +"</span></a>"
    td_comments.innerHTML = "<button id='openComments_"+tx.Txid +"' onclick='openComments(`"+ tx.Txid +"`)'><img class='comment' id='comment_"+ tx.Txid +"' height='20' src='icons/"+ commentImage +"'><span class='commentcounter' id='comment_counter_"+ tx.Txid +"'>"+ tx.Comments +"</span></button>"
  }else{
    commentImage = "comment_0.png"
    // td_comments.innerHTML = "<a title='comment' onclick='reply(`"+ tx.Txid +"`,`"+ testcomment +"`)'><img class='comment' id='comment_"+ tx.Txid +"' height='20' src='icons/"+ commentImage +"'><span class='commentcounter' id='comment_counter_"+ tx.Txid +"'>"+ tx.Comments +"</span></a>"
    td_comments.innerHTML = "<button id='openComments_"+tx.Txid +"' onclick='openComments(`"+ tx.Txid +"`)'><img class='comment' id='comment_"+ tx.Txid +"' height='20' src='icons/"+ commentImage +"'><span class='commentcounter' id='comment_counter_"+ tx.Txid +"'>"+ tx.Comments +"</span></button>"
  }

  td_like.innerHTML = "<a title='like' onclick='like(`"+ tx.Txid +"`)'><img class='like' id='like_"+ tx.Txid +"' height='20' src='icons/"+ likeImage +"'><span class='likecounter' id='like_counter_"+ tx.Txid +"'>"+ tx.Likes +"</span></a>"
  td_sender.innerHTML = tx.txid
  td_blockheight.innerHTML = (confirmed) ? (tx.BlockHeight) : ("unconfirmed")
  td_score.innerHTML = tx.Score

  var link = check_link(tx.Link)
  var type = check_type(tx.DataType)
  var title = check_title(tx.Title)

  if (link && type && title){
    td_6a_magnethash.innerHTML = "<a class='' href='"+ link +"'><img height='15' src='icons/icons8-magnet-filled-50.png'>" + link + "</a>";
    td_6a_title.innerHTML = title;
    td_6a_type.innerHTML = type;

    input_data = '"' + link + '","' + title + '","' + tx.Txid + '"'
    td_play.innerHTML = "<button title='play with webtorrent' class='result-play' onclick='play(" + input_data + ");'><span class='glyphicon glyphicon-play-circle'></span></button>";
    td_play.style.width = "15px";


    tr.appendChild(td_txid);
    tr.appendChild(td_like);
    tr.appendChild(td_comments);
    tr.appendChild(td_score);
    tr.appendChild(td_play);
    tr.appendChild(td_6a_title);
    tr.appendChild(td_6a_magnethash);
    tr.appendChild(td_blockheight);

    document.getElementById('bitdb_output').appendChild(tr);
  }
};

function get_video_data(hash,title,sender){
  // Insert Title
  document.getElementById('video_title').innerHTML = title

  var query = {
    request: {
      encoding: {
        b1: "hex"
      },
      find: {
        b1: { "$in": ["e902"] },
        s2: {
          "$regex": hash, "$options": "i"
        },
        'senders.a' :  {
          "$in": [sender]
        }
      },
      project: {
        b0:1 ,b1: 1, s2: 1, tx: 1, block_index: 1, _id: 0, senders: 1
      },
      limit: 10
    },
    response: {
      encoding: {
        b1: "hex"
      }
    }
  };
  var b64 = btoa(JSON.stringify(query));
  var url = "https://bitdb.network/v2/q/" + b64;

  var header = {
    headers: { key: "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44" }
  };

  fetch(url, header).then(function(r) {
    return r.json()
  }).then(function(r) {

    for(i in r['confirmed']){
      var tx = r['confirmed'][i];
      console.log(tx.s2);
      var p = document.createElement('p');
      p.innerHTML = tx.s2
      document.getElementById('video_description').appendChild(p)
    };

  })
}

function download_torrent(hash,title,sender){
  var torrentId = hash + "&tr=udp://explodie.org:6969&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.empire-js.us:1337&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://tracker.opentrackr.org:1337&tr=wss://tracker.openwebtorrent.com"
  // var torrentId = "magnet:?xt=urn:btih:" + hash + "&tr=udp://explodie.org:6969&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.empire-js.us:1337&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://tracker.opentrackr.org:1337&tr=wss://tracker.openwebtorrent.com"


  var client = new WebTorrent()

  // HTML elements
  var $body = document.body
  var $progressBar = document.querySelector('#progressBar')
  var $numPeers = document.querySelector('#numPeers')
  var $downloaded = document.querySelector('#downloaded')
  var $total = document.querySelector('#total')
  var $remaining = document.querySelector('#remaining')
  var $uploadSpeed = document.querySelector('#uploadSpeed')
  var $downloadSpeed = document.querySelector('#downloadSpeed')

  // Download the torrent
  client.add(torrentId, function (torrent) {

    // insert data
    document.getElementById('torrentLink').innerHTML = torrentId
    // Video Data from Blockchain
    get_video_data(hash,title,sender);

    // show divs
    document.getElementById('video_output_container').style.display = "block";

    // Torrents can contain many files. Let's use the .mp4 file
    var file = torrent.files.find(function (file) {
      return file.name.endsWith('.mp4')
    })

    // Stream the file in the browser
    file.appendTo('#output')

    // Trigger statistics refresh
    torrent.on('done', onDone)
    setInterval(onProgress, 500)
    onProgress()

    // Statistics
    function onProgress () {
      // Peers
      $numPeers.innerHTML = torrent.numPeers + (torrent.numPeers === 1 ? ' peer' : ' peers')

      // Progress
      var percent = Math.round(torrent.progress * 100 * 100) / 100
      $progressBar.style.width = percent + '%'
      $downloaded.innerHTML = prettyBytes(torrent.downloaded)
      $total.innerHTML = prettyBytes(torrent.length)

      // Remaining time
      var remaining
      if (torrent.done) {
        remaining = 'Done.'
      } else {
        remaining = moment.duration(torrent.timeRemaining / 1000, 'seconds').humanize()
        remaining = remaining[0].toUpperCase() + remaining.substring(1) + ' remaining.'
      }
      $remaining.innerHTML = remaining

      // Speed rates
      $downloadSpeed.innerHTML = prettyBytes(torrent.downloadSpeed) + '/s'
      $uploadSpeed.innerHTML = prettyBytes(torrent.uploadSpeed) + '/s'
    }
    function onDone () {
      $body.className += ' is-seed'
      onProgress()
    }
  })
}

// Human readable bytes util
function prettyBytes(num) {
  var exponent, unit, neg = num < 0, units = ['B', 'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  if (neg) num = -num
  if (num < 1) return (neg ? '-' : '') + num + ' B'
  exponent = Math.min(Math.floor(Math.log(num) / Math.log(1000)), units.length - 1)
  num = Number((num / Math.pow(1000, exponent)).toFixed(2))
  unit = units[exponent]
  return (neg ? '-' : '') + num + ' ' + unit
}
