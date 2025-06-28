async function logar() {
    const nickname = document.getElementById("nickname").value;
    const password = document.getElementById("password").value;

    if(!nickname || !password){
        alert("preencha todos os dados");
        return
    }


    const res = await fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type":  "application/json"
        },
        body: JSON.stringify({ nickname, password })
    });


    if(!res.ok){
        alert("Login falhou");
        return;
    }

    const data = await res.json();
    localStorage.setItem("token", data.token);
    localStorage.setItem("nickname", data.nickname);

    window.location.href = "/game.html"



    
    
    

}