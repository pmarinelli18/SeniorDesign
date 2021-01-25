import React, { useEffect, useState, PureComponent } from "react";
import { StyleSheet, AppRegistry, View } from "react-native";
import { Container, Header, Content, Button, Text } from "native-base";
import { Actions } from "react-native-router-flux";

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		alignItems: "center",
		backgroundColor: "transparent",
	},
});
export default class Start extends PureComponent {
	constructor() {
		super();
	}

	render() {
		return (
			<View {...this.props} style={styles.container}>
				<Button onPress={() => Actions._login()}>
					<Text>Log in</Text>
				</Button>
			</View>
		);
	}
}

// const styles = StyleSheet.create({
// 	container: {
// 		flex: 1,
// 		backgroundColor: "#FFF",
// 	},
// });

AppRegistry.registerComponent("Start", () => BestGameEver);
{
	/* <Text>Replace screen</Text>
			<Text>Prop from dynamic method</Text>
			<Button onPress={() => Actions.Login()}>Back</Button> */
}
