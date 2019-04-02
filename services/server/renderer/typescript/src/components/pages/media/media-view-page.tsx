import * as React from "react";
import "./media-view-page.scss"
import {IGlobalProps} from "../../../global-props";
import {IMediaResponse} from "./media-list-page";
import {Video} from "../../media/video";

interface IProps extends IGlobalProps {
    view: IMediaResponse;
}

export class MediaViewPage extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        return <div id={"media-view-page"}>
            <div className={"row"}>
                <div className={"col-sm-8 offset-sm-2"}>
                    <Video media={this.props.view}/>
                </div>
            </div>
        </div>;
    }

}
