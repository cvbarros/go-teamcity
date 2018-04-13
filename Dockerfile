FROM jetbrains/teamcity-server:2017.2.1

COPY reset.sh /reset.sh
CMD ["/reset.sh"]