import React, { useEffect, useState, PureComponent } from "react";
import { StyleSheet, AppRegistry, View } from "react-native";
import { Container, Header, Content, Button, Text } from "native-base";
import { Actions } from "react-native-router-flux";
export default class Login extends PureComponent {
	constructor() {
		super();
	}

	render() {
		return (
			<View {...this.props} style={styles.container}>
				<Button onPress={() => Actions._start()}>
					<Text>Test</Text>
				</Button>
			</View>
		);
	}
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#FFF",
	},
});

AppRegistry.registerComponent("Login", () => BestGameEver);
