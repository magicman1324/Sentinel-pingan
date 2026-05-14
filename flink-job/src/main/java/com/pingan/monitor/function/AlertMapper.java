package com.pingan.monitor.function;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.pingan.monitor.model.AlertEvent;
import org.apache.flink.api.common.functions.MapFunction;

public class AlertMapper implements MapFunction<AlertEvent, String> {

    private static final ObjectMapper mapper = new ObjectMapper();

    @Override
    public String map(AlertEvent alert) throws Exception {
        return mapper.writeValueAsString(alert);
    }
}
