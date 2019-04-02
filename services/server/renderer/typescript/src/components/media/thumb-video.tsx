import * as React from "react";
import "./thumb-video.scss"
import {IMediaResponse} from "../pages/media/media-list-page";
import {Video} from "./video";

interface IProps {
    media: IMediaResponse;
    isLink: boolean;
    showMetadata: boolean;
}

interface IDimensions {
    width: number;
    height: number;
}

const MAX_AXIS_SIZE: number = 200;

export class ThumbVideo extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        if (this.props.media.thumbnails.length === 0) {
            return null;
        }

        if (this.props.isLink) {
            return <a className={"thumb-video-link"} href={`/media/view/${this.props.media.id}`}>
                {this.renderMainView()}
            </a>
        }

        return this.renderMainView();
    }

    private renderMainView(): JSX.Element {
        const dim: IDimensions = this.getDimensions();
        const style: React.CSSProperties = {
            width: "100%",
            minHeight: `${dim.height}px`,
            minWidth: `${dim.width}px`
        };

        return <div className={"thumb-video"}>
            <div className={"media-thumbnail-video"}>
                <img className="media-thumbnail-video-thumbnail" src={this.props.media.thumbnails[0]} style={style}/>
                <Video media={this.props.media} style={style}/>
            </div>
            {this.renderMetadata()}
        </div>;
    }

    private renderMetadata(): JSX.Element | null {
        if (!this.props.showMetadata) {
            return null;
        }

        const desc: string = this.props.media.description.length > 20 ? `${this.props.media.description.substring(0, 20)}...` : this.props.media.description;

        return <div className={"metadata"}>
            <span className={"description"}>{desc + " " + this.props.media.hashtags.join(" ")}</span>
        </div>;
    }

    private getDimensions(): IDimensions {
        let height: number = this.props.media.videos[0].height;
        let width: number = this.props.media.videos[0].width;
        const ratio: number = width / MAX_AXIS_SIZE;
        height = height / ratio;

        return {
            height: height,
            width: MAX_AXIS_SIZE
        }
    }

}