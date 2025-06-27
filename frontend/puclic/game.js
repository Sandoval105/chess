
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



const game = new Chess()

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



function onDrop(source, target) {

    const move = game.move({
        from: source,
        to: target,
        promotion: "q"
    });

    if (move == null)  return 'snapback';



   updateStatus()

};


const board = Chessboard("board", {
    draggable: true,
    position: 'start',
    pieceTheme: function(piece){
        return 'pieces/' + (customPieceMap[piece] || piece);
    },
    onDrop: onDrop

});







