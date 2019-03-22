import {IPageProps, Page} from "../page";
import * as React from "react";
import {registerComponent} from "preact-rpc";

interface IProps extends IPageProps {

}

export class HomePage extends React.Component<IProps, {}> {
    
    public render(): JSX.Element {
        return <Page id={"home-page"} authPayload={this.props.authPayload}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    Hello!
                </div>
            </div>
        </Page>;
    }
    
}

registerComponent('HomePage', HomePage);
