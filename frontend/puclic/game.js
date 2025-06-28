
let board ;
const game = new Chess()

const nickname = localStorage.getItem("nickname");
const ws = new WebSocket("ws://localhost:8080/ws?nickname=" + nickname);

let color = "white";





const customPieceMap = {
    'wP': 'Bpeao.svg',
    'wR': 'Btorre.svg',
    'wN': 'Bcavalo.svg',
    'wB': 'Bbispo.svg',
    'wQ': 'Brainha.svg',
    'wK': 'Brei.svg',
    'bP': 'Ppeao.svg',
    'bR': 'Ptorre.svg',
    'bN': 'Pcavalo.svg',
    'bB': 'Pbispo.svg',
    'bQ': 'Prainha.svg',
    'bK': 'Prei.svg',
};





function updateStatus(){
    let status = "";

    if(game.in_checkmate()){
        status = "check mate"
    }else if(game.in_draw()){
        status = "rei afogado"
    }else if(game.in_check()){
        status = "check"
    } else{
        status = "seu turno"
    
    }

    document.getElementById("status").innerText = status;

}



ws.onmessage = (msg) => {
    try {
        const data =  JSON.parse(msg.data);

        if(data.type == "start"){
            color = data.color;
           
            board = Chessboard("board", {
                draggable: true,
                position: 'start',
                pieceTheme: function(piece){
                    return 'pieces/' + (customPieceMap[piece] || piece);
                },
                onDrop: onDrop,
                orientation: color
            });
            

            updateStatus();

            return
        }
         
        if (data.type == "time"){

            console.log("seu time: ", data.your);
            console.log("oponente time: ",data.enemy)

            return
        }

        if (data.from && data.to){
            game.move({from: data.from, to: data.to });
            board.position(game.fen());
            updateStatus();
        }

    } catch (e) {
        console.log(e)
    }


}




function onDrop(source, target) {

    if (game.turn() !== color[0]) return 'snapback';
    const move = game.move({
        from: source,
        to: target,
        promotion: "q"
    });

    if (move == null)  return 'snapback';
    


   updateStatus()

   ///enviar jogadas pro oponente 

   ws.send(JSON.stringify({
    type: "move",
    from: source,
    to: target
   }));

};





