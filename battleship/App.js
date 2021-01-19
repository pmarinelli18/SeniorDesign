import React, { PureComponent } from "react";
import { StyleSheet, Text, View, AppRegistry, StatusBar } from "react-native";
import { registerRootComponent } from "expo";
import { Finger } from "./renderers";
import { GameEngine } from "react-native-game-engine";
// import Game from "./game";
export default class BestGameEver extends PureComponent {
	constructor() {
		super();
	}

	render() {
		return (
			<GameEngine
				style={{ backgroundColor: "black" }}
				entities={{
					1: { position: [40, 200], renderer: <Finger /> },
					2: { position: [100, 200], renderer: <Finger /> },
					3: { position: [160, 200], renderer: <Finger /> },
					4: { position: [220, 200], renderer: <Finger /> },
					5: { position: [280, 200], renderer: <Finger /> },
				}}
			>
				<Text style={[styles.setFontSize, styles.setColorRed]}>Hi</Text>
				<StatusBar hidden={true} />
			</GameEngine>
		);
	}
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#FFF",
	},
});

AppRegistry.registerComponent("BestGameEver", () => BestGameEver);

// export default function App() {
//   return (
//     <View style={styles.container}>
//       <Text>Open up App.js to start working on your app!</Text>
//       <StatusBar style="auto" />
//     </View>
//   );
// }

// const styles = StyleSheet.create({
// 	container: {
// 		flex: 1,
// 		backgroundColor: "#fff",
// 		alignItems: "center",
// 		justifyContent: "center",
// 	},
// });
