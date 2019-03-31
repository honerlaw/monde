import * as React from "react";
import {IMediaResponse} from "./upload-list-page";
import {ThumbVideo} from "../media/thumb-video";

interface IProps {
    media: IMediaResponse[];
}

export class HomePage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"home-page"}>
            <div className={"row"}>
                {this.renderVideos()}
            </div>
        </div>;
    }

    private renderVideos(): JSX.Element[] {
        let elements: JSX.Element[] = [];

        for (let i: number = 0; i < 30; ++i) {
            const media: IMediaResponse = this.props.media[0];
            elements.push(<div className={"col col-sm-6"}><ThumbVideo isLink={true} key={media.id + i} media={media}/></div>);
        }

        return elements;
    }

}
