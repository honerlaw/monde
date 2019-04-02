import * as React from "react";
import "./not-found-page.scss";

export class NotFoundPage extends React.Component<{}, {}> {

    public render(): JSX.Element {
        return <div id={"not-found-page"}>
            <img src={require("../../../../../../assets/img/404.svg")} />
            <h3>Page Not Found!</h3>
        </div>;
    }

}
