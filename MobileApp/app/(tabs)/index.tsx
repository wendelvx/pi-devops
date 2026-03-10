import { View, Text, TouchableOpacity, TextInput, StyleSheet } from "react-native";
import { useState, useEffect } from "react";
import { io } from "socket.io-client";

const socket = io("http://SEU_IP:3000"); // coloque o IP do seu servidor

export default function HomeScreen() {

  const [nickname, setNickname] = useState("");
  const [classeSelecionada, setClasseSelecionada] = useState<string | null>(null);

  const classes = [
    "Front-end",
    "Back-end",
    "DevOps",
    "QA",
    "Security"
  ];

  function selecionarClasse(classe:string) {

    if (classeSelecionada) return;

    setClasseSelecionada(classe);

  }

  function atacar() {

    if (!nickname) return;

    const payload = {
      nickname: nickname
    };

    socket.emit("attack", payload);

  }

  return (
    <View style={styles.container}>

      <Text style={styles.title}>SQUAD FORMANDO</Text>

      <TextInput
        style={styles.input}
        placeholder="Digite seu nickname"
        placeholderTextColor="#777"
        value={nickname}
        onChangeText={setNickname}
      />

      <Text style={styles.subtitle}>Escolha sua classe</Text>

      {classes.map((classe) => {

        const selecionado = classeSelecionada === classe;

        return (
          <TouchableOpacity
            key={classe}
            onPress={() => selecionarClasse(classe)}
            style={[
              styles.card,
              selecionado && styles.cardSelecionado
            ]}
          >

            <Text style={[
              styles.text,
              selecionado && styles.textSelecionado
            ]}>
              {classe}
            </Text>

          </TouchableOpacity>
        );

      })}

      {classeSelecionada && (
        <Text style={styles.role}>
          Você é: {classeSelecionada}
        </Text>
      )}

      <TouchableOpacity style={styles.attackButton} onPress={atacar}>
        <Text style={styles.attackText}>ATACAR ⚔️</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({

  container: {
    flex: 1,
    backgroundColor: "#000",
    padding: 25,
    justifyContent: "center"
  },

  title: {
    fontSize: 30,
    color: "#fff",
    fontWeight: "bold",
    marginBottom: 20
  },

  subtitle: {
    color: "#888",
    marginBottom: 10
  },

  input: {
    borderColor: "#333",
    borderWidth: 1,
    borderRadius: 10,
    padding: 12,
    color: "white",
    marginBottom: 20
  },

  card: {
    borderColor: "#333",
    borderWidth: 1,
    padding: 15,
    borderRadius: 10,
    marginBottom: 10
  },

  cardSelecionado: {
    borderColor: "#00ff66"
  },

  text: {
    color: "#555"
  },

  textSelecionado: {
    color: "#00ff66",
    fontWeight: "bold"
  },

  role: {
    color: "#888",
    textAlign: "center",
    marginTop: 10
  },

  attackButton: {
    backgroundColor: "#ff0033",
    padding: 20,
    borderRadius: 10,
    alignItems: "center",
    marginTop: 25
  },

  attackText: {
    color: "white",
    fontSize: 20,
    fontWeight: "bold"
  },
});