import * as React from "react";
import {registerComponent} from "preact-rpc";

interface IProps {
    toWhat: string;
}

export class Index extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div>Hello {this.props.toWhat} </div>;
    }

}

registerComponent('Index', Index);
