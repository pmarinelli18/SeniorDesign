import React, { useEffect, useState, PureComponent } from "react";
import { Text, StyleSheet, AppRegistry } from "react-native";
export default class Login extends PureComponent {
	constructor() {
		super();
	}

	render() {
		return <Text> Hi</Text>;
	}
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#FFF",
	},
});

AppRegistry.registerComponent("Login", () => BestGameEver);
