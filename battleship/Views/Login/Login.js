import React, { useEffect, useState, PureComponent } from "react";
import { StyleSheet, AppRegistry, View } from "react-native";
import {
	Container,
	Header,
	Content,
	Button,
	Text,
	Form,
	Item,
	Input,
} from "native-base";
import { Actions } from "react-native-router-flux";

const Login = (props) => {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");

	const handleUsernameChange = (e) => {
		setUsername(e);
		// console.log(e);
	};

	const handlePasswordChange = (e) => {
		setPassword(e);
		// console.log(e);
	};

	const handleSubmit = (e) => {
		// comment out e.preventDefault(); for the Ant Design Form
		//e.preventDefault();
		//sourced from https://medium.com/@maison.moa/setting-up-an-express-backend-server-for-create-react-app-bc7620b20a61
		console.log("This is the Username: ");
		console.log(username);
		console.log("This is the Password: ");
		console.log(password);
		Actions._home();
		// 	const login = async () => {
		// 		const response = await axios.post("/login", {
		// 			email,
		// 			password,
		// 		});
		// 		const body = response.data;

		// 		if (response.status !== 200) {
		// 			throw Error(body.error);
		// 		}

		// 		console.log(body);
		// 		Actions._home();
		// 	};

		// 	login().catch((err) => {
		// 		alert("Wrong email or password");
		// 		console.log(err);
		// 	});
	};

	// if (props.authorized) Actions._start();

	return (
		<View style={styles.container}>
			<Form>
				<Item regular style={styles.formItem}>
					<Input
						placeholder="Username"
						onChangeText={(text) => handleUsernameChange(text)}
					/>
				</Item>
				<Item regular style={styles.formItem}>
					<Input
						placeholder="Password"
						onChangeText={(text) => handlePasswordChange(text)}
					/>
				</Item>
			</Form>
			<Button
				onPress={() => {
					handleSubmit();
				}}
				style={styles.formItem}
			>
				<Text>Login</Text>
			</Button>
		</View>
	);
};

// export default class Login extends PureComponent {
// 	constructor() {
// 		super();
// 	}

// 	render() {
// 		return (
// 			<View {...this.props} style={styles.container}>
// 				<Form>
// 					<Item regular style={styles.formItem}>
// 						<Input placeholder="Username" />
// 					</Item>
// 					<Item regular style={styles.formItem}>
// 						<Input placeholder="Password" />
// 					</Item>
// 				</Form>
// 				<Button onPress={() => Actions._home()} style={styles.formItem}>
// 					<Text>Login</Text>
// 				</Button>
// 			</View>
// 		);
// 	}
// }

const styles = StyleSheet.create({
	container: {
		flex: 1,
		padding: 24,
		justifyContent: "center",
		backgroundColor: "transparent",
	},
	formItem: {
		margin: 5,
	},
});

export default Login;
// AppRegistry.registerComponent("Login", () => BestGameEver);
