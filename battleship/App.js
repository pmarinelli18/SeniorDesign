import React from "react";
import { Platform, StyleSheet, Text } from "react-native";
// import { StackViewStyleInterpolator } from "react-navigation-stack"
import {
	Scene,
	Router,
	Actions,
	ActionConst,
	Overlay,
	Tabs,
	Modal,
	Drawer,
	Stack,
	Lightbox,
} from "react-native-router-flux";
// import Home from "./Views/Home";
import Start from "./Views/Start/Start";
import Login from "./Views/Login/Login";
import Home from "./Views/Home/Home";

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "transparent",
		justifyContent: "center",
		alignItems: "center",
	},
	scene: {
		backgroundColor: "#F5FCFF",
		shadowOpacity: 1,
		shadowRadius: 3,
	},
	tabBarStyle: {
		backgroundColor: "#eee",
	},
	tabBarSelectedItemStyle: {
		backgroundColor: "#ddd",
	},
});

const stateHandler = (prevState, newState, action) => {
	console.log("onStateChange: ACTION:", action);
};

// on Android, the URI prefix typically contains a host in addition to scheme
const prefix = Platform.OS === "android" ? "mychat://mychat/" : "mychat://";

// const transitionConfig = () => ({
// 	screenInterpolator: StackViewStyleInterpolator.forFadeFromBottomAndroid,
// });

const App = () => (
	<Router
		onStateChange={stateHandler}
		sceneStyle={styles.scene}
		uriPrefix={prefix}
	>
		<Overlay key="overlay">
			<Modal key="modal" hideNavBar>
				<Scene key="_start" component={Start} title="Start" />
				<Scene key="_login" component={Login} title="Login" />
				<Scene key="_home" component={Home} title="Home" />
			</Modal>
		</Overlay>
	</Router>
);
export default App;
