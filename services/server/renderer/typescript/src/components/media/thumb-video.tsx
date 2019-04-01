import * as React from "react";
import {IMediaResponse, IMediaVideoResponse} from "../pages/upload-list-page";
import "./thumb-video.scss"

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
        const mp4: string | null = this.getMp4Url();
        if (this.props.media.thumbnails.length === 0 || !mp4) {
            return null;
        }

        if (this.props.isLink) {
            return <a className={"thumb-video-link"} href={`/media/view/${this.props.media.id}`}>
                {this.renderMainView(mp4)}
            </a>
        }

        return this.renderMainView(mp4);
    }

    private renderMainView(mp4: string): JSX.Element {
        const dim: IDimensions = this.getDimensions();
        const style: React.CSSProperties = {
            width: "100%",
            minHeight: `${dim.height}px`,
            minWidth: `${dim.width}px`
        };

        return <div className={"thumb-video"}>
            <div className={"media-thumbnail-video"}>
                <img className="media-thumbnail-video-thumbnail" src={this.props.media.thumbnails[0]} style={style}/>
                <video controls={false} autoPlay={true} loop={true} style={style}>
                    <source src={mp4} type="video/mp4"/>
                </video>
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

    private getMp4Url(): string | null {
        const videos: IMediaVideoResponse[] = this.props.media.videos.filter((video) => video.type === "mp4");
        if (videos.length > 0) {
            return videos[0].url;
        }
        return null;
    }

}