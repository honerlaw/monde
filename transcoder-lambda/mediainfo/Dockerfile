FROM amazonlinux:2017.03

# Install Development Tools necessary to compile MediaInfo
RUN yum -y groupinstall 'Development Tools'
# Install library required to add CURL support to Mediainfo
RUN yum -y install libcurl-devel wget

# Download MediaInfo
RUN wget https://mediaarea.net/download/binary/mediainfo/18.12/MediaInfo_CLI_18.12_GNU_FromSource.tar.xz
# Untar MediaInfo
RUN tar xvf MediaInfo_CLI_18.12_GNU_FromSource.tar.xz

WORKDIR /MediaInfo_CLI_GNU_FromSource

# build with libcurl so we can pass a url
RUN ./CLI_Compile.sh --with-libcurl

WORKDIR /MediaInfo_CLI_GNU_FromSource/MediaInfo/Project/GNU/CLI
